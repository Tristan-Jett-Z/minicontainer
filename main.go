package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		child()
		return
	}
	
	parent()
}

func parent(){
	if len(os.Args)<2{
		panic("用法:minicontainer <command> [args]")
	}

	args:=append([]string{"child"},os.Args[1:]...)

	cmd := exec.Command("/proc/self/exe", 
			args...,
		)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("启动容器进程中...")

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	
	cgroupPath:=setupCgroup(cmd.Process.Pid)

	defer os.RemoveAll(cgroupPath)

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func child() {
	fmt.Println("我是容器子进程")
        fmt.Println("PID:", os.Getpid())
	
	if err:=syscall.Mount("","/","",
			syscall.MS_PRIVATE|
			syscall.MS_REC,
			"",
	);err!=nil{
		panic(err)
	}

	rootfs := "./rootfs"
	if err := syscall.Chroot(rootfs); err != nil {
		panic(err)
	}

	if err := syscall.Chdir("/"); err != nil {
		panic(err)
	}

	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		panic(err)
	}

	if len(os.Args)<3{
		panic("没有指定容器内要执行的命令")
	}

	command:=os.Args[2]

	path,err:=exec.LookPath(command)
	if err!=nil{
		panic(err)
	}
	
	err=syscall.Exec(
		path,
		os.Args[2:],
		os.Environ(),
	)

	if err != nil {
		panic(err)
	}
}

func setupCgroup(pid int) string{
	cgroupPath := "/sys/fs/cgroup/minicontainer"
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		cgroupPath+"/memory.max",
		[]byte("100M"),
		0644,
	); err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		cgroupPath+"/cgroup.procs",
		[]byte(fmt.Sprintf("%d", pid)),
		0644,
	); err != nil {
		panic(err)
	}
	return cgroupPath
}
