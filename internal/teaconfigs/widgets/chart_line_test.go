package widgets

import "testing"

func TestLineChart_AllParamNames(t *testing.T) {
	{
		chart := &LineChart{}

		{
			line := NewLine()
			line.Param = "${param1}"
			chart.AddLine(line)
		}

		{
			line := NewLine()
			line.Param = "${param2} * 10"
			chart.AddLine(line)
		}

		{
			line := NewLine()
			line.Param = "${param2.value} * ${param3}"
			chart.AddLine(line)
		}

		{
			line := NewLine()
			line.Param = "${ param2 } * ${param3}"
			chart.AddLine(line)
		}

		{
			line := NewLine()
			line.Param = "param4"
			chart.AddLine(line)
		}

		t.Log(chart.AllParamNames())
	}

	{
		chart := &LineChart{}
		t.Log(chart.AllParamNames())
	}
}
