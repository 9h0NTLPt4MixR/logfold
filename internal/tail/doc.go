// Package tail implements a lightweight log-file tailer for logfold.
//
// It opens a file, seeks to its current end, and emits each newly appended
// line over a buffered channel.  The caller drives the lifecycle via a
// [context.Context]; cancelling the context causes [Tailer.Run] to return
// [context.Canceled].
//
// Typical usage:
//
//	tr := tail.New("/var/log/app.log")
//	go func() {
//		if err := tr.Run(ctx); err != nil && err != context.Canceled {
//			log.Printf("tailer error: %v", err)
//		}
//	}()
//	for line := range tr.Lines {
//		// hand off to parser / aggregator pipeline
//	}
//
// The tailer polls the file at a fixed [PollInterval] when the OS reports
// EOF, making it suitable for plain files as well as named pipes.
package tail
