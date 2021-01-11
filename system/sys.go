package system

/*
#include <unistd.h>
int GoFork() {
	return fork();
}
*/
import "C"

func Fork() int {
	return C.GoFork()
}
