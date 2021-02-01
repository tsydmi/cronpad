<template>
  <canvas :id="elementId"></canvas>
</template>

<script>
import Chart from "chart.js";

export default {
  name: "PieChart",

  props: {
    elementId: {
      type: String,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    value: {
      type: Object,
      required: false,
    },
    sum: {
      type: Number,
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
        type: 'pie',
        data: this.value,
        options: {
          legend: {
            position: 'bottom'
          },
          title: {
            display: true,
            text: this.title
          },
          tooltips: {
            callbacks: {
              title: function (tooltipItems, data) {
                let value = data.datasets[0].data[tooltipItems[0].index]

                if (value === 1 || value === '1') {
                  return ` ${value} hour`
                } else {
                  return ` ${value} hours`
                }
              }.bind(this),
              label: function (tooltipItems, data) {
                if (this.sum) {
                  let value = data.datasets[0].data[tooltipItems.index]
                  return ` ${(value*100/this.sum).toFixed(2)} %`
                }

                return ''
              }.bind(this),
            },
          },
        },
      });
    },
    updateChart() {
      this.chart.data.datasets = this.value.datasets
      this.chart.data.labels = this.value.labels
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