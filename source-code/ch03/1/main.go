func main() {  
    x := 1
    fmt.Println(x)     // prints 1
    {
        fmt.Println(x) // prints 1
        x := 2
        fmt.Println(x) // prints 2
    }
    fmt.Println(x)     // prints 1 (不是2)
}