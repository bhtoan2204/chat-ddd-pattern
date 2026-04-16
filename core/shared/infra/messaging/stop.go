package messaging

import "sync"

func StopConsumers(consumers []Consumer) {
	var wg sync.WaitGroup

	for _, consumer := range consumers {
		if consumer == nil {
			continue
		}

		wg.Add(1)
		go func(consumer Consumer) {
			defer wg.Done()
			consumer.Stop()
		}(consumer)
	}

	wg.Wait()
}
