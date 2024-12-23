
options (
    go_package="github.com/basecomplextech/baselibrary/proto/pclock"
)

// HLTimestamp is a hybrid logical timestamp.
struct HLTimestamp {
    Wall    int64;  // physical time in unix nanoseconds
    Seq     uint32; // logical sequence number
}
