package devops

import (
	"fmt"
	"github.com/timescale/tsbs/pkg/data"
	"github.com/timescale/tsbs/pkg/data/usecases/common"
	"math/rand"
	"time"
)

var (
	labelNet             = []byte("net") // heap optimization
	labelNetTagInterface = []byte("interface")
	labelNetTagIp        = []byte("ip-address")

	// Reuse NormalDistributions as arguments to other distributions. This is
	// safe to do because the higher-level distribution advances the ND and
	// immediately uses its value and saves the state
	highND = common.ND(50, 1)
	lowND  = common.ND(5, 1)

	netFields = []common.LabeledDistributionMaker{
		{Label: []byte("bytes_sent"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 50) }},
		{Label: []byte("bytes_recv"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 50) }},
		{Label: []byte("packets_sent"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 50) }},
		{Label: []byte("packets_recv"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 50) }},
		{Label: []byte("err_in"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 5) }},
		{Label: []byte("err_out"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 5) }},
		{Label: []byte("drop_in"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 5) }},
		{Label: []byte("drop_out"), DistributionMaker: func() common.Distribution { return common.WD(common.ND(0, 1), 5) }},
	}
)

type NetMeasurement struct {
	*common.SubsystemMeasurement
	interfaceName string
    ipAddress string
}

func randomIp() string {
    return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func NewNetMeasurement(start time.Time) *NetMeasurement {
	sub := common.NewSubsystemMeasurementWithDistributionMakers(start, netFields)
	interfaceName := fmt.Sprintf("eth%d", rand.Intn(4))
    ipAddress := randomIp()
	return &NetMeasurement{
		SubsystemMeasurement: sub,
		interfaceName:        interfaceName,
        ipAddress:            ipAddress,
	}
}

func (m *NetMeasurement) ToPoint(p *data.Point) {
	m.ToPointAllInt64(p, labelNet, netFields)
	p.AppendTag(labelNetTagInterface, m.interfaceName)
	p.AppendTag(labelNetTagIp, m.ipAddress)
    if rand.Float64() < 0.2 {
        m.ipAddress = randomIp()
    }
}
