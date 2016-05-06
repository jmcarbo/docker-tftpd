package main

import (
	"github.com/pin/tftp"
	"os"
	"io"
	"fmt"
	"time"
)

func writeHanlder(filename string, w io.WriterTo) error {
    fmt.Printf("Receiving file %s\n", filename)
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    // In case client provides tsize option.
    if t, ok := w.(tftp.IncomingTransfer); ok {
        if n, ok := t.Size(); ok {
            fmt.Printf("Transfer size: %d\n", n)
        }
    }
    n, err := w.WriteTo(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    fmt.Printf("%d bytes received\n", n)
    return nil
}

func readHandler(filename string, r io.ReaderFrom) error {
    fmt.Printf("Sending file %s\n", filename)
    file, err := os.Open(filename)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    // Optional tsize support.
    // Set transfer size before calling ReadFrom.
    if t, ok := r.(tftp.OutgoingTransfer); ok {
        if fi, err := file.Stat(); err == nil {
            t.SetSize(fi.Size())
        }
    }
    n, err := r.ReadFrom(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        return err
    }
    fmt.Printf("%d bytes sent\n", n)
    return nil
}

func main() {
    // use nil in place of handler to disable read or write operations
    s := tftp.NewServer(readHandler, writeHanlder)
    s.SetTimeout(25 * time.Second) // optional
    fmt.Println("Starting tftpd server")
    err := s.ListenAndServe(":69") // blocks until s.Shutdown() is called
    if err != nil {
        fmt.Fprintf(os.Stdout, "server: %v\n", err)
        os.Exit(1)
    }
}
