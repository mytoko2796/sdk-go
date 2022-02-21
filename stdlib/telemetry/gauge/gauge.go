package gauge

import (
	"runtime"
	"time"

	"github.com/mytoko2796/sdk-go/stdlib/telemetry/exporter"
	"go.opencensus.io/metric"
	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricproducer"
)

const (
	DimensionLess  = metricdata.UnitDimensionless
	DimensionBytes = metricdata.UnitBytes
	DimensionMSec  = metricdata.UnitMilliseconds
)

var (
	MRegistry = metric.NewRegistry()

	//cpu
	GaugeCPUCount, _  = MRegistry.AddInt64Gauge(`proc_cpu_count`, metric.WithDescription(`Number of machine CPU`), metric.WithUnit(DimensionLess))
	GaugeGoroutine, _ = MRegistry.AddInt64Gauge(`proc_num_goroutine`, metric.WithDescription(`Number of machine CPU`), metric.WithUnit(DimensionLess))
	GaugeCgoCall, _   = MRegistry.AddInt64Gauge(`proc_cgo_call`, metric.WithDescription(`Sum of Cgo calls`), metric.WithUnit(DimensionLess))

	//general
	GaugeMemAlloc, _      = MRegistry.AddInt64Gauge(`proc_mem_alloc`, metric.WithDescription(`Alloc is bytes of allocated heap objects.`), metric.WithUnit(DimensionBytes))
	GaugeMemTotalAlloc, _ = MRegistry.AddInt64Gauge(`proc_mem_alloctotal`, metric.WithDescription(`TotalAlloc is cumulative bytes allocated for heap objects.`), metric.WithUnit(DimensionBytes))
	GaugeMemSys, _        = MRegistry.AddInt64Gauge(`proc_mem_sys`, metric.WithDescription(`Sys is the total bytes of memory obtained from the OS.`), metric.WithUnit(DimensionBytes))
	GaugeMemLookUps, _    = MRegistry.AddInt64Gauge(`proc_mem_lookups`, metric.WithDescription(`Lookups is the number of pointer lookups performed by the runtime.`), metric.WithUnit(DimensionLess))
	GaugeMemMallocs, _    = MRegistry.AddInt64Gauge(`proc_mem_mallocs`, metric.WithDescription(`Mallocs is the cumulative count of heap objects allocated`), metric.WithUnit(DimensionLess))
	GaugeMemFrees, _      = MRegistry.AddInt64Gauge(`proc_mem_frees`, metric.WithDescription(`Frees is the cumulative count of heap objects freed.`), metric.WithUnit(DimensionLess))

	//heap
	GaugeMemHeapAlloc, _    = MRegistry.AddInt64Gauge(`proc_memheap_alloc`, metric.WithDescription(`HeapAlloc is bytes of allocated heap objects.`), metric.WithUnit(DimensionBytes))
	GaugeMemHeapSys, _      = MRegistry.AddInt64Gauge(`proc_memheap_sys`, metric.WithDescription(`HeapSys is bytes of heap memory obtained from the OS.`), metric.WithUnit(DimensionBytes))
	GaugeMemHeapIdle, _     = MRegistry.AddInt64Gauge(`proc_memheap_idle`, metric.WithDescription(`HeapIdle is bytes in idle (unused) spans.`), metric.WithUnit(DimensionBytes))
	GaugeMemHeapInUse, _    = MRegistry.AddInt64Gauge(`proc_memheap_inuse`, metric.WithDescription(`HeapInuse is bytes in in-use spans.`), metric.WithUnit(DimensionBytes))
	GaugeMemHeapReleased, _ = MRegistry.AddInt64Gauge(`proc_memheap_released`, metric.WithDescription(`HeapReleased is bytes of physical memory returned to the OS.`), metric.WithUnit(DimensionBytes))
	GaugeMemHeapObj, _      = MRegistry.AddInt64Gauge(`proc_memheap_objects`, metric.WithDescription(`HeapObjects is the number of allocated heap objects.`), metric.WithUnit(DimensionLess))

	// Stack
	GaugeMemStackInUse, _      = MRegistry.AddInt64Gauge(`proc_memstack_inuse`, metric.WithDescription(`StackInuse is bytes in stack spans.`), metric.WithUnit(DimensionBytes))
	GaugeMemStackSys, _        = MRegistry.AddInt64Gauge(`proc_memstack_sys`, metric.WithDescription(`StackSys is bytes of stack memory obtained from the OS.`), metric.WithUnit(DimensionBytes))
	GaugeMemStackSpanInUse, _  = MRegistry.AddInt64Gauge(`proc_memstack_mspaninuse`, metric.WithDescription(`MSpanInuse is bytes of allocated mspan structures.`), metric.WithUnit(DimensionBytes))
	GaugeMemStackSpanSys, _    = MRegistry.AddInt64Gauge(`proc_memstack_mspansys`, metric.WithDescription(`MSpanSys is bytes of memory obtained from the OS for mspan structures.`), metric.WithUnit(DimensionBytes))
	GaugeMemStackCacheInUse, _ = MRegistry.AddInt64Gauge(`proc_memstack_mcacheinuse`, metric.WithDescription(`MCacheInuse is bytes of allocated mcache structures.`), metric.WithUnit(DimensionBytes))
	GaugeMemStackCacheSys, _   = MRegistry.AddInt64Gauge(`proc_memstack_cachesys`, metric.WithDescription(`MCacheSys is bytes of memory obtained from the OS for mcache structures.`), metric.WithUnit(DimensionBytes))

	//others
	GaugeMemOtherSys, _ = MRegistry.AddInt64Gauge(`proc_memother_sys`, metric.WithDescription(`OtherSys is bytes of memory in miscellaneous off-heap runtime allocations.`), metric.WithUnit(DimensionBytes))

	//GC
	GaugeGCSys, _         = MRegistry.AddInt64Gauge(`proc_gc_sys`, metric.WithDescription(`GCSys is bytes of memory in garbage collection metadata.`), metric.WithUnit(DimensionBytes))
	GaugeGCNext, _        = MRegistry.AddInt64Gauge(`proc_gc_next`, metric.WithDescription(`NextGC is the target heap size of the next GC cycle.`), metric.WithUnit(DimensionBytes))
	GaugeGCLast, _        = MRegistry.AddInt64Gauge(`proc_gc_last`, metric.WithDescription(`Seconds from last GC`), metric.WithUnit(DimensionMSec))
	GaugeGCPause, _       = MRegistry.AddInt64Gauge(`proc_gc_pause`, metric.WithDescription(`Last GC pause`), metric.WithUnit(DimensionMSec))
	GaugeGCPauseTotal, _  = MRegistry.AddInt64Gauge(`proc_gc_pausetotal`, metric.WithDescription(`PauseTotalNs is the cumulative nanoseconds in GC stop-the-world pauses since the program started.`), metric.WithUnit(DimensionMSec))
	GaugeGCCycleCount, _  = MRegistry.AddInt64Gauge(`proc_gc_cyclecount`, metric.WithDescription(`NumGC is the number of completed GC cycles.`), metric.WithUnit(DimensionLess))
	GaugeGCCPUFraction, _ = MRegistry.AddFloat64Gauge(`proc_gc_cpufrac`, metric.WithDescription(`GCCPUFraction is the fraction of this program's available CPU time used by the GC since the program started.`), metric.WithUnit(DimensionLess))
)

func Init(conf exporter.StatsOptions, termSig chan struct{}) error {
	metricproducer.GlobalManager().AddProducer(MRegistry)
	recordGauges(conf.RecordPeriod, termSig)
	return nil
}
func recordGauges(period time.Duration, termSig chan struct{}) error {
	var (
		rtm runtime.MemStats
	)
	GaugeCPUCountEntry, err := GaugeCPUCount.GetEntry()
	if err != nil {
		return err
	}
	GaugeGoroutineEntry, err := GaugeGoroutine.GetEntry()
	if err != nil {
		return err
	}
	GaugeCgoCallEntry, err := GaugeCgoCall.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemAllocEntry, err := GaugeMemAlloc.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemTotalAllocEntry, err := GaugeMemTotalAlloc.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemSysEntry, err := GaugeMemSys.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemLookUpsEntry, err := GaugeMemLookUps.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemMallocsEntry, err := GaugeMemMallocs.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemFreesEntry, err := GaugeMemFrees.GetEntry()
	if err != nil {
		return err
	}

	//Heap
	GaugeMemHeapAllocEntry, err := GaugeMemHeapAlloc.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemHeapSysEntry, err := GaugeMemHeapSys.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemHeapIdleEntry, err := GaugeMemHeapIdle.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemHeapInUseEntry, err := GaugeMemHeapInUse.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemHeapReleasedEntry, err := GaugeMemHeapReleased.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemHeapObjEntry, err := GaugeMemHeapObj.GetEntry()
	if err != nil {
		return err
	}

	//Stack
	GaugeMemStackInUseEntry, err := GaugeMemStackInUse.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemStackSysEntry, err := GaugeMemStackSys.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemStackSpanInUseEntry, err := GaugeMemStackSpanInUse.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemStackSpanSysEntry, err := GaugeMemStackSpanSys.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemStackCacheInUseEntry, err := GaugeMemStackCacheInUse.GetEntry()
	if err != nil {
		return err
	}
	GaugeMemStackCacheSysEntry, err := GaugeMemStackCacheSys.GetEntry()
	if err != nil {
		return err
	}

	//others
	GaugeMemOtherSysEntry, err := GaugeMemOtherSys.GetEntry()
	if err != nil {
		return err
	}

	//GC
	GaugeGCSysEntry, err := GaugeGCSys.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCNextEntry, err := GaugeGCNext.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCLastEntry, err := GaugeGCLast.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCPauseTotalEntry, err := GaugeGCPauseTotal.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCPauseEntry, err := GaugeGCPause.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCCycleCountEntry, err := GaugeGCCycleCount.GetEntry()
	if err != nil {
		return err
	}
	GaugeGCCPUFractionEntry, err := GaugeGCCPUFraction.GetEntry()
	if err != nil {
		return err
	}
	tick := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-termSig:
				return
			case <-tick.C:
				runtime.ReadMemStats(&rtm)
				//cpu
				GaugeCPUCountEntry.Set(int64(runtime.NumCPU()))
				GaugeGoroutineEntry.Set(int64(runtime.NumGoroutine()))
				GaugeCgoCallEntry.Set(runtime.NumCgoCall())

				//general
				GaugeMemAllocEntry.Set(int64(rtm.Alloc))
				GaugeMemTotalAllocEntry.Set(int64(rtm.TotalAlloc))
				GaugeMemSysEntry.Set(int64(rtm.Sys))
				GaugeMemLookUpsEntry.Set(int64(rtm.Lookups))
				GaugeMemMallocsEntry.Set(int64(rtm.Mallocs))
				GaugeMemFreesEntry.Set(int64(rtm.Frees))

				//Heap
				GaugeMemHeapAllocEntry.Set(int64(rtm.HeapAlloc))
				GaugeMemHeapSysEntry.Set(int64(rtm.HeapSys))
				GaugeMemHeapIdleEntry.Set(int64(rtm.HeapIdle))
				GaugeMemHeapInUseEntry.Set(int64(rtm.HeapInuse))
				GaugeMemHeapReleasedEntry.Set(int64(rtm.HeapReleased))
				GaugeMemHeapObjEntry.Set(int64(rtm.HeapObjects))

				//Stack
				GaugeMemStackInUseEntry.Set(int64(rtm.StackInuse))
				GaugeMemStackSysEntry.Set(int64(rtm.StackSys))
				GaugeMemStackSpanInUseEntry.Set(int64(rtm.MSpanInuse))
				GaugeMemStackSpanSysEntry.Set(int64(rtm.MSpanSys))
				GaugeMemStackCacheInUseEntry.Set(int64(rtm.MCacheInuse))
				GaugeMemStackCacheSysEntry.Set(int64(rtm.MCacheSys))

				//others
				GaugeMemOtherSysEntry.Set(int64(rtm.OtherSys))
				//GC
				GaugeGCSysEntry.Set(int64(rtm.GCSys))
				GaugeGCNextEntry.Set(int64(rtm.NextGC))
				GaugeGCLastEntry.Set(int64(rtm.LastGC))
				GaugeGCPauseTotalEntry.Set(int64(rtm.PauseTotalNs))
				GaugeGCPauseEntry.Set(int64(rtm.PauseNs[(rtm.NumGC+255)%256]))
				GaugeGCCycleCountEntry.Set(int64(rtm.NumGC))

				GaugeGCCPUFractionEntry.Set(rtm.GCCPUFraction)
			}
		}
	}()
	return nil
}
