///|
fn main {
  let args = @sys.get_cli_args()
  let program = "flow"
  if args.length() <= 1 {
    print_usage(program)
    return
  }
  let command = args[1]
  if command == "help" || command == "--help" || command == "-h" {
    print_usage(program)
  } else if command == "version" || command == "--version" {
    print_version()
  } else if command == "greet" {
    handle_greet(args)
  } else if command == "fib" {
    handle_fib(program, args)
  } else {
    println("\{program}: unknown command `\{command}`")
    println("run `\{program} help` to see available commands")
    @sys.exit(1)
  }
}

///|
fn print_usage(program : String) -> Unit {
  println("Usage: \{program} <command> [arguments]")
  println("")
  println("Commands:")
  println("  help                   Show this message")
  println("  version                Print the CLI version")
  println("  greet [name]           Print a friendly greeting")
  println("  fib <n>                Compute nth Fibonacci number using lib.fib")
}

///|
fn print_version() -> Unit {
  println("flow 0.1.0")
}

///|
fn handle_greet(args : Array[String]) -> Unit {
  let name = if args.length() > 2 { args[2] } else { "friend" }
  println("Hello, \{name}!")
}

///|
fn handle_fib(program : String, args : Array[String]) -> Unit {
  if args.length() <= 2 {
    println("\{program}: `fib` requires an integer argument")
    println("example: \{program} fib 10")
    @sys.exit(1)
  }
  let n = parse_integer_or_exit(program, args[2])
  if n < 0 {
    println("\{program}: fibonacci index must be non-negative")
    @sys.exit(1)
  }
  let result = @lib.fib(n)
  println("fib(\{n}) = \{result}")
}

///|
fn parse_integer_or_exit(program : String, raw : String) -> Int {
  if raw.length() == 0 {
    println("\{program}: expected an integer but found empty input")
    @sys.exit(1)
  }
  let zero = '0'.to_int()
  let nine = '9'.to_int()
  let mut index = 0
  let mut negative = false
  if raw[0] == '-'.to_int() {
    negative = true
    index = 1
  }
  if index >= raw.length() {
    println("\{program}: expected digits after '-' in integer argument")
    @sys.exit(1)
  }
  let mut value = 0
  while index < raw.length() {
    let code = raw[index]
    if code < zero || code > nine {
      println("\{program}: not a valid integer: \{raw}")
      @sys.exit(1)
    }
    value = value * 10 + (code - zero)
    index = index + 1
  }
  if negative {
    -value
  } else {
    value
  }
}
