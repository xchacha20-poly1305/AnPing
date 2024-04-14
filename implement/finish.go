package implement

func ListenOnDone(done <-chan struct{}, callback func()) {
	<-done
	callback()
}
