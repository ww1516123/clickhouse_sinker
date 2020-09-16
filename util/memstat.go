package util
// import (
//     "runtime"
//     "syscall"
//     "log"
// )
//
// type MemStatus struct {
//     All  uint64 `json:"all"`
//     Used uint64 `json:"used"`
//     Free uint64 `json:"free"`
//     Self uint64 `json:"self"`
// }
//
// func MemStat() MemStatus {
//     //自身占用
//     memStat := new(runtime.MemStats)
//     runtime.ReadMemStats(memStat)
//     mem := MemStatus{}
//     mem.Self = memStat.Alloc
//
//     //系统占用,仅linux/mac下有效
//     //system memory usage
//     sysInfo := new(syscall.Sysinfo_t)
//     err := syscall.Sysinfo(sysInfo)
//     if err == nil {
//         mem.All = sysInfo.Totalram
//         mem.Free = sysInfo.Freeram
//         mem.Used = mem.All - mem.Free
//     }
//     memPrint(mem)
//     return mem
// }
//
// func memPrint(memStat MemStatus){
//   log.Printf("Used:  %dMb", memStat.Used/1024/1024.0)
//   log.Printf("Total:      %dMb", memStat.All/1024/1024.0)
//   log.Printf("Free:     %dMb", memStat.Free/1024/1024.0)
//   log.Printf("Self:     %dMb", memStat.Self/1024/1024.0)
// }
