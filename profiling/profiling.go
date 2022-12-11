package profiling

import (
	"bufio"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	"github.com/complex1tech/baselibrary/errors2"
)

type Profiling struct {
	config *Config
	mu     sync.Mutex

	cpuFile   *os.File
	memFile   *os.File
	blockFile *os.File
	mutexFile *os.File
}

// New copies the config and returns a new profiling instance.
func New(config *Config) *Profiling {
	config1 := &Config{}
	*config1 = *config
	return &Profiling{config: config1}
}

// Init reads a profiling config from a json file if it exists and returns a new profiling instance.
func Init(path string) (*Profiling, error) {
	config, err := ReadConfig(path)
	switch {
	case os.IsNotExist(err):
		config = Default()
	case err != nil:
		return nil, err
	}

	p := New(config)
	return p, nil
}

// Start starts enabled profilers.
func (p *Profiling) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.config.Enabled {
		return nil
	}
	ok := false

	// cpu
	if err := p.startCPU(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopCPU()
	}()

	// memory
	if err := p.startMemory(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopMemory()
	}()

	// block
	if err := p.startBlock(); err != nil {
		return err
	}
	defer func() {
		if ok {
			return
		}
		p.stopBlock()
	}()

	// mutex
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
func (p *Profiling) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err0 := p.stopCPU()
	err1 := p.stopMemory()
	err2 := p.stopBlock()
	err3 := p.stopMutex()

	return errors2.Combine(err0, err1, err2, err3)
}

// cpu

func (p *Profiling) startCPU() error {
	cpu := p.config.CPU
	if !cpu.Enabled {
		return nil
	}

	// create file
	f, err := os.Create(cpu.Path)
	if err != nil {
		return err
	}

	// close on error
	ok := false
	defer func() {
		if ok {
			return
		}
		f.Close()
	}()

	// set profile rate
	if cpu.Rate > 0 {
		runtime.SetCPUProfileRate(cpu.Rate)
	}

	// start profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}

	// done
	ok = true
	p.cpuFile = f
	return nil
}

func (p *Profiling) stopCPU() error {
	if p.cpuFile == nil {
		return nil
	}
	defer func() {
		p.cpuFile.Close()
		p.cpuFile = nil
	}()

	// stop profiling
	pprof.StopCPUProfile()

	// close file
	if err := p.cpuFile.Close(); err != nil {
		return err
	}
	return nil
}

// memory

func (p *Profiling) startMemory() error {
	mem := p.config.Memory
	if !mem.Enabled {
		return nil
	}

	// create file
	f, err := os.Create(mem.Path)
	if err != nil {
		return err
	}

	// set profile rate
	if mem.Rate > 0 {
		runtime.MemProfileRate = mem.Rate
	}

	p.memFile = f
	return nil
}

func (p *Profiling) stopMemory() error {
	if p.memFile == nil {
		return nil
	}
	defer func() {
		p.memFile.Close()
		p.memFile = nil
	}()

	// get up-to-date stats
	runtime.GC()

	// write memory profile
	w := bufio.NewWriter(p.memFile)
	if err := pprof.WriteHeapProfile(w); err != nil {
		return err
	}

	// flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.memFile.Close(); err != nil {
		return err
	}
	return nil
}

// block

func (p *Profiling) startBlock() error {
	block := p.config.Block
	if !block.Enabled {
		return nil
	}

	// create file
	f, err := os.Create(block.Path)
	if err != nil {
		return err
	}

	// set profile rate
	rate := block.Rate
	if rate == 0 {
		rate = DefaultBlockProfileRate
	}
	runtime.SetBlockProfileRate(rate)

	p.blockFile = f
	return nil
}

func (p *Profiling) stopBlock() (err error) {
	if p.blockFile == nil {
		return nil
	}
	defer func() {
		p.blockFile.Close()
		p.blockFile = nil
	}()

	// write block profile
	w := bufio.NewWriter(p.blockFile)
	if err := pprof.Lookup("block").WriteTo(w, 0); err != nil {
		return err
	}

	// flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.blockFile.Close(); err != nil {
		return err
	}
	return nil
}

// mutex

func (p *Profiling) startMutex() error {
	mutex := p.config.Mutex
	if !mutex.Enabled {
		return nil
	}

	// create file
	f, err := os.Create(mutex.Path)
	if err != nil {
		return err
	}

	// set profile rate
	rate := mutex.Rate
	if rate == 0 {
		rate = DefalutMutexProfileRate
	}
	runtime.SetMutexProfileFraction(rate)

	p.mutexFile = f
	return nil
}

func (p *Profiling) stopMutex() (err error) {
	if p.mutexFile == nil {
		return nil
	}
	defer func() {
		p.mutexFile.Close()
		p.mutexFile = nil
	}()

	// write mutex profile
	w := bufio.NewWriter(p.mutexFile)
	if err := pprof.Lookup("mutex").WriteTo(w, 0); err != nil {
		return err
	}

	// flush and close file
	if err := w.Flush(); err != nil {
		return err
	}
	if err := p.mutexFile.Close(); err != nil {
		return err
	}
	return nil
}
