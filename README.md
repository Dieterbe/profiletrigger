automatically trigger a profile in your golang (go) application when a condition is matched.

[godoc](http://godoc.org/github.com/Dieterbe/profiletrigger)

# currently implemented:

* when certain number of bytes allocated, save a heap (memory) profile
* when cpu usage reaches a certain percentage, save a cpu profile.

# demo

see the included cpudemo and heapdemo programs, which gradually add more and cpu and heap utilisation, to show the profiletrigger kicking in.
