# check-intel-microcode-version
On Linux and Darwin, attempts to match cpu and microcode against list of those receiving the [microcode-20191115 update](https://github.com/intel/Intel-Linux-Processor-Microcode-Data-Files/releases/tag/microcode-20191115).

This code's not incredibly robust, but it's not very complex either.  It expects to be able to run "sysctl" on Darwin and to run "egrep" on "/proc/cpuinfo" on Linux.  The CPU matching is imperfect because neither of those platforms supplies all the necessary information about their processor (as far as I can tell).

Sample run:
```
go run .
This processor is {processor: stepping: family:6 model:8e steppingID:9 platformID: oldVersion:202 newVersion:202 products:}
This computer matches {processor:AML-Y22 stepping:H0 family:6 model:8e steppingID:9 platformID:10 oldVersion:198 newVersion:202 products:Core Gen8 Mobile}, this microcode is 202
This computer matches {processor:KBL-U/Y stepping:H0 family:6 model:8e steppingID:9 platformID:c0 oldVersion:198 newVersion:202 products:Core Gen7 Mobile}, this microcode is 202
This computer matches {processor:KBL-U23e stepping:J1 family:6 model:8e steppingID:9 platformID:c0 oldVersion:198 newVersion:202 products:Core Gen7 Mobile}, this microcode is 202
```

That is, this processor's microcode is 202, and that is the updated version for the all three processors that match what we know about this processor.
