package rxgo

import "time"

func throttleLast(in chan interface{}, delay time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(delay)
		var last interface{}

		for {
			select {
			case <-ticker.C:
				out <- last

			case value, ok := <-in:
				if !ok {
					continue
				}
				last = value
			}
		}
	}()

	return out
}

func throttleFirst(in chan interface{}, delay time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(delay)
		var last interface{}

		last = <-in
		out <- last

		for {
			select {
			case <-ticker.C:
				out <- last

			case value, ok := <-in:
				if !ok {
					continue
				}
				last = value
			}
		}
	}()

	return out
}

func debounce(in chan interface{}, delay time.Duration) chan interface{} {
	out := make(chan interface{})

	go func() {
		ticker := time.NewTicker(delay)
		var last interface{}

		for {
			select {
			case <-ticker.C:
				out <- last

			case value, ok := <-in:
				if !ok {
					continue
				}
				last = value
				ticker.Stop()
				ticker = time.NewTicker(delay)
			}
		}
	}()

	return out
}

func buffer(in chan interface{}, delay time.Duration) chan []interface{} {
	out := make(chan []interface{})

	go func() {
		ticker := time.NewTicker(delay)
		buffer := make([]interface{}, 10)

		for {
			select {
			case <-ticker.C:
				out <- buffer
				buffer = buffer[:0]

			case value, ok := <-in:
				if !ok {
					continue
				}
				buffer = append(buffer, value)
			}
		}
	}()

	return out
}

func window(in chan interface{}, delay time.Duration) chan chan interface{} {
	out := make(chan chan interface{})

	go func() {
		ticker := time.NewTicker(delay)
		buffer := make([]interface{}, 10)

		for {
			select {
			case <-ticker.C:
				ch := make(chan interface{}, len(buffer))
				for _, value := range buffer {
					ch <- value
				}
				out <- ch
				buffer = buffer[:0]

			case value, ok := <-in:
				if !ok {
					continue
				}
				buffer = append(buffer, value)
			}
		}
	}()

	return out
}

func windowSize(in chan interface{}, size int) chan chan interface{} {
	out := make(chan chan interface{})

	go func() {
		buffer := make([]interface{}, 10)

		for {
			select {
			case value, ok := <-in:
				if !ok {
					continue
				}
				buffer = append(buffer, value)
				if len(buffer) < size {
					continue
				}

				ch := make(chan interface{}, len(buffer))
				for _, value := range buffer {
					ch <- value
				}
				out <- ch
				buffer = buffer[:0]
			}
		}
	}()

	return out
}

// func xthrottleLastX(in chan interface{}) chan interface{} {
// 	out := make(chan interface{}, 0)
// 	var last interface{}

// 	go func() {
// 		for {
// 			select {
// 			case value, ok := <-in:
// 				if !ok {
// 					return
// 				}
// 				last = value

// 			case out <- last:
// 			}
// 		}
// 	}()
// 	return out
// }

// func xthrottleFirstX(in chan interface{}) chan interface{} {
// 	out := make(chan interface{}, 0)
// 	var last interface{}

// 	go func() {
// 		for {
// 			select {
// 			case out <- last:
// 			default:
// 				select {
// 				case value, ok := <-in:
// 					if !ok {
// 						return
// 					}
// 					last = value
// 				}
// 			}
// 		}
// 	}()
// 	return out
// }

// func debounce(in chan interface{}, delay time.Duration) chan interface{} {
// 	ticker := time.NewTicker(delay)
// 	out := make(chan interface{}, 0)
// 	var last interface{}

// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				select {
// 				case out <- last:
// 				default:
// 				}
// 			default:
// 				select {
// 				case value, ok := <-in:
// 					if !ok {
// 						return
// 					}
// 					last = value
// 				}

// 			}
// 		}
// 	}()
// 	return out
// }
