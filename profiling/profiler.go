// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package profiling

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
)

type profiler struct {
	config *Config
	mu     sync.Mutex

	cpuFile   *os.File
	memFile   *os.File
	blockFile *os.File
	mutexFile *os.File
}

// newProfiler copies the config and returns a new profiling instance.
func newProfiler(config *Config) *profiler {
	config1 := &Config{}
	*config1 = *config
	return &profiler{config: config1}
}

// Start starts enabled profilers.
func (p *profiler) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.config.Enabled {
		return nil
	}
	ok := false

	// Cpu
	if err := p.startCPU(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopCPU()
	}()

	// Memory
	if err := p.startMemory(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopMemory()
	}()

	// Block
	if err := p.startBlock(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopBlock()
	}()

	// Mutex
	if err := p.startMutex(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopMutex()
	}()

	ok = true
	return nil
}

// Stop stops all profilers.
func (p *profiler) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err0 := p.stopCPU()
	err1 := p.stopMemory()
	err2 := p.stopBlock()
	err3 := p.stopMutex()

	return errors.Join(err0, err1, err2, err3)
}

// cpu

func (p *profiler) startCPU() error {
	cpu := p.config.CPU
	if !cpu.Enabled {
		return nil
	}

	// Create dir
	if err := makeParentDir(cpu.Path); err != nil {
		return err
	}

	// Create file
	f, err := os.Create(cpu.Path)
	if err != nil {
		return err
	}

	// Close on error
	ok := false
	defer func() {
		if ok {
			return
		}
		f.Close()
	}()

	// Set profile rate
	if cpu.Rate > 0 {
		runtime.SetCPUProfileRate(cpu.Rate)
	}

	// Start profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}

	// Done
	ok = true
	p.cpuFile = f
	return nil
}

func (p *profiler) stopCPU() error {
	if p.cpuFile == nil {
		return nil
	}
	defer func() {
		p.cpuFile.Close()
		p.cpuFile = nil
	}()

	// Stop profiling
	pprof.StopCPUProfile()

	// Close file
	if err := p.cpuFile.Close(); err != nil {
		return err
	}
	return nil
}

// memory

func (p *profiler) startMemory() error {
	mem := p.config.Memory
	if !mem.Enabled {
		return nil
	}

	// Create dir
	if err := makeParentDir(mem.Path); err != nil {
		return err
	}

	// Create file
	f, err := os.Create(mem.Path)
	if err != nil {
		return err
	}

	// Set profile rate
	if mem.Rate > 0 {
		runtime.MemProfileRate = mem.Rate
	}

	p.memFile = f
	return nil
}

func (p *profiler) stopMemory() error {
	if p.memFile == nil {
		return nil
	}
	defer func() {
		p.memFile.Close()
		p.memFile = nil
	}()

	// Get up-to-date stats
	runtime.GC()

	// Write memory profile
	w := bufio.NewWriter(p.memFile)
	if err := pprof.WriteHeapProfile(w); err != nil {
		return err
	}

	// Flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.memFile.Close(); err != nil {
		return err
	}
	return nil
}

// block

func (p *profiler) startBlock() error {
	block := p.config.Block
	if !block.Enabled {
		return nil
	}

	// Create dir
	if err := makeParentDir(block.Path); err != nil {
		return err
	}

	// Create file
	f, err := os.Create(block.Path)
	if err != nil {
		return err
	}

	// Set profile rate
	rate := block.Rate
	if rate == 0 {
		rate = DefaultBlockProfileRate
	}
	runtime.SetBlockProfileRate(rate)

	p.blockFile = f
	return nil
}

func (p *profiler) stopBlock() (err error) {
	if p.blockFile == nil {
		return nil
	}
	defer func() {
		p.blockFile.Close()
		p.blockFile = nil
	}()

	// Write block profile
	w := bufio.NewWriter(p.blockFile)
	if err := pprof.Lookup("block").WriteTo(w, 0); err != nil {
		return err
	}

	// Flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.blockFile.Close(); err != nil {
		return err
	}
	return nil
}

// mutex

func (p *profiler) startMutex() error {
	mutex := p.config.Mutex
	if !mutex.Enabled {
		return nil
	}

	// Create dir
	if err := makeParentDir(mutex.Path); err != nil {
		return err
	}

	// Create file
	f, err := os.Create(mutex.Path)
	if err != nil {
		return err
	}

	// Set profile rate
	rate := mutex.Rate
	if rate == 0 {
		rate = DefalutMutexProfileRate
	}
	runtime.SetMutexProfileFraction(rate)

	p.mutexFile = f
	return nil
}

func (p *profiler) stopMutex() (err error) {
	if p.mutexFile == nil {
		return nil
	}
	defer func() {
		p.mutexFile.Close()
		p.mutexFile = nil
	}()

	// Write mutex profile
	w := bufio.NewWriter(p.mutexFile)
	if err := pprof.Lookup("mutex").WriteTo(w, 0); err != nil {
		return err
	}

	// Flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.mutexFile.Close(); err != nil {
		return err
	}
	return nil
}

// makedir

const dirMode = 0755

func makeParentDir(file string) error {
	path := filepath.Dir(file)

	_, err := os.Stat(path)
	switch {
	case err == nil:
		return nil
	case !errors.Is(err, os.ErrNotExist):
		return err
	}

	return os.MkdirAll(path, dirMode)
}
