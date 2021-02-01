<template>
  <canvas :id="elementId"></canvas>
</template>

<script>
import Chart from "chart.js";

const BACKGROUND_COLORS = ["#FFB300", "#E0E0E0"]

export default {
  name: "ProgressChart",

  props: {
    elementId: {
      type: String,
      required: true,
    },
    value: {
      type: Object,
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
      const ctx = document.getElementById(this.elementId)

      return new Chart(ctx, {
        type: 'doughnut',
        data: {
          datasets: this.getDatasets,
          labels: ['done', 'ahead'],
        },
        options: {
          rotation: Math.PI,
          circumference: Math.PI,
          legend: {
            display: false,
          },
          title: {
            display: true,
            text: this.getTitle,
            position: 'bottom',
            fontSize: 20,
          },
          tooltips: {
            enabled: true,
            mode: 'single',
            callbacks: {
              label: function (tooltipItems, data) {
                let value = data.datasets[0].data[tooltipItems.index]

                if (value === 1 || value === '1') {
                  return ` ${value} day`
                }

                return ` ${value} days`
              }
            },
          },
        },
      });
    },
    updateChart() {
      this.chart.data.datasets = this.getDatasets
      this.chart.options.title.text = this.getTitle
      this.chart.update();
    },
  },
  computed: {
    getTitle() {
      if (this.value && this.value.daysAhead != null && this.value.daysPassed != null) {
        if (this.value.daysPassed === 0 && this.value.daysPassed === 0) {
          return 'Progress unknown'
        }

        return `${((100 * this.value.daysPassed) / (this.value.daysPassed + this.value.daysAhead)).toFixed(2)} %`;
      }

      return 'No Data'
    },
    getDatasets() {
      return [{
        data: [this.value.daysPassed, this.value.daysAhead],
        backgroundColor: BACKGROUND_COLORS,
      }]
    }
  },
  watch: {
    'value': function () {
      this.createOrUpdateChart()
    }
  }
}
</script>