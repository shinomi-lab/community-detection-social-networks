package optimization

// diffuse.goを使ってモンテカルロ法で影響関数を求めている

import (
	diff "difftools/diffusion"
	"difftools/network"
	"math/rand"
	"runtime"
	"sync"
)

// 乱数シード値の代わりに乱数生成器`r`を渡す
//
// return: result of mont (配列で影響関数の答えがinfoごとにある)
func RunInflProp(
	sampleSize int,
	net network.Network,
	// adj [][]int,
	seedSet []diff.SeedInfo,
	userProbTable diff.UserProbTable,
	pop diff.PopList,
	interestList diff.InterestList,
	assumList diff.AssumList,
	r *rand.Rand,
) map[diff.InfoType]float64 {

	numWorkers := runtime.NumCPU()
	if sampleSize < numWorkers {
		numWorkers = sampleSize
	}
	// 各ワーカーに割り当てる試行回数
	samplesPerWorker := sampleSize / numWorkers
	extraSamples := sampleSize % numWorkers

	// 1. 各ワーカー用に独立した乱数シードを生成（データレース回避）
	workerSeeds := make([]int64, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerSeeds[i] = r.Int63()
	}

	type partialResult struct {
		sumF int
		sumT int
	}

	resultsChan := make(chan partialResult, numWorkers)
	var wg sync.WaitGroup

	// 2. CPUコア数分ワーカーを起動して試行を分散
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		count := samplesPerWorker
		if i == 0 {
			count += extraSamples // 余り分を最初のワーカーに追加
		}
		seed := workerSeeds[i]

		go func(count int, seed int64) {
			defer wg.Done()

			// ワーカー専用の乱数生成器を作成
			localR := rand.New(rand.NewSource(seed))

			localSumF := 0
			localSumT := 0

			for k := 0; k < count; k++ {
				dist := diff.Diffuse(net, seedSet, userProbTable, pop, interestList, assumList, localR)
				localSumF += len(dist[diff.InfoType_F])
				localSumT += len(dist[diff.InfoType_T])
			}

			resultsChan <- partialResult{sumF: localSumF, sumT: localSumT}
		}(count, seed)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 3. 各ワーカーの集計結果を加算
	totalSumF := 0
	totalSumT := 0
	for res := range resultsChan {
		totalSumF += res.sumF
		totalSumT += res.sumT
	}

	// 4. sampleSize で割って平均値を算出
	ans := make(map[diff.InfoType]float64, diff.InfoTypes_n)
	ans[diff.InfoType_F] = float64(totalSumF) / float64(sampleSize)
	ans[diff.InfoType_T] = float64(totalSumT) / float64(sampleSize)
	/*
		// n := len(adj)

		// dist := make([][]int, diff.InfoTypes_n)
		ans := make(map[diff.InfoType]float64, diff.InfoTypes_n)

		for i := 0; i < sampleSize; i++ {
			var dist = diff.Diffuse(net, seedSet, userProbTable, pop, interestList, assumList, r) //-1 is correct?
			ans[diff.InfoType_F] += float64(len(dist[diff.InfoType_F])) / float64(sampleSize)
			ans[diff.InfoType_T] += float64(len(dist[diff.InfoType_T])) / float64(sampleSize)
		}
	*/
	return ans
}
