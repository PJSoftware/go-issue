package util

import "log/slog"

func PanicOn(err error, a ...any) {
	if err != nil {
		x := []any{"ERROR", err.Error()}
		if len(a)%2 == 1 {
			x = append(x, "context")
		}
		x = append(x, a...)

		slog.Error("panic triggered:", x...)
		panic(err)
	}
}
