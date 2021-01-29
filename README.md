# lambda-debugger-go
Time travel debugging directly from AWS Lambda with the Go runtime

# =(

After a couple of months I have realized that this library (for the time being) is not feasible.
Lambda does not allow `ptrace` system call (obviously) for security reason.
