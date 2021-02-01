<template>
  <canvas :id="elementId"></canvas>
</template>

<script>
import Chart from "chart.js"
import dayjs from "dayjs"

export default {
  name: "DateRangeBarChart",

  props: {
    elementId: {
      type: String,
      required: true,
    },
    value: {
      type: Array,
      required: false,
    },
  },
  data: () => ({
    chart: null,
  }),
  methods: {
    createOrUpdateChart() {
      if (this.chart) {
        this.updateChart()
      } else {
        this.chart = this.createChart()
      }
    },
    createChart() {
      const ctx = document.getElementById(this.elementId);
      return new Chart(ctx, {
        type: 'bar',
        data: {
          datasets: this.value,
        },
        options: {
          scales: {
            xAxes: [{
              offset: true,
              type: 'time',
              time: {
                unit: 'day',
              },
              stacked: true,
            }],
            yAxes: [{
              stacked: true,
            }],
          },
          legend: {
            display: false,
            labels: {
              defaultFontFamily: 'Varela Round'
            },
          },
          layout: {
            padding: {
              right: 20
            },
          },
          tooltips: {
            enabled: true,
            mode: 'single',
            callbacks: {
              title: function(tooltipItem) {
                return dayjs(tooltipItem[0].label).format('YYYY-MM-DD')
              },
            },
          },
        },
      })
    },
    updateChart() {
      this.chart.data.datasets = this.value
      this.chart.update();
    },
  },
  watch: {
    'value': function () {
      this.createOrUpdateChart()
    }
  }
}
</script>