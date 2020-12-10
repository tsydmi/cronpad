<template>
  <div>
    <v-calendar
        class="pa-5"
        weekdays="1, 2, 3, 4, 5, 6, 0"
        :now="today"
        :value="value"
    >
      <template v-slot:day-label="{day, weekday, month, year, present, date}">
        <div class="" @click="selectDay(date)">
          <div
              v-bind:class="['d-flex justify-center',
              {
                'v-calendar-month--selected-week': isSelectedWeek(date),
                'secondary--text': !isCurrentMonth(month),
                'primary--text font-weight-bold': present
              }]"
          >
            <span v-if="day === 1" class="pr-1">
              {{ monthNames[month - 1] }}
            </span>
            <span>
              {{ day }}
            </span>
          </div>
        </div>
      </template>

      <template v-slot:day="{ past, date }">
        <v-row
            class="align-end mr-1 ml-1"
        >
          <template v-if="dayHasEvents(date)">
            <v-sheet
                v-for="(event, i) in getEvents(date)"
                :key="i"
                :title="event.name"
                :color="event.color"
                :width="`${100*getEventDurationHours(event)/8}%`"
                height="5px"
                tile
            >
            </v-sheet>
          </template>
        </v-row>
      </template>
    </v-calendar>
  </div>
</template>

<script>
export default {
  props: {
    today: {
      type: String,
      required: true,
    },
    value: {
      type: String,
      required: true,
    },
    events: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    selectedWeekFirstDay: null,
    selectedWeekLastDay: null,
    selectedMonth: null,
    monthNames: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
  }),
  methods: {
    selectDay(date) {
      this.updateFirstAndLastDayOfSelection(date)
      this.$emit('changeCalendarValue', date)
    },
    isSelectedWeek(date) {
      let selectedDate = new Date(date);
      return selectedDate >= this.selectedWeekFirstDay && selectedDate <= this.selectedWeekLastDay
    },
    updateFirstAndLastDayOfSelection(date) {
      let selectedDate = new Date(date)
      this.selectedMonth = selectedDate.getMonth()

      let realWeekDay = (selectedDate.getDay() + 6) % 7

      this.selectedWeekFirstDay = new Date(selectedDate.setDate(selectedDate.getDate() - realWeekDay))
      this.selectedWeekLastDay = new Date(selectedDate.setDate(selectedDate.getDate() + 6))
    },
    isCurrentMonth(month) {
      return this.selectedMonth === (month - 1)
    },
    dayHasEvents(date) {
      var start = new Date(date);
      start.setHours(0, 0, 0, 0);

      var end = new Date(date);
      end.setHours(23, 59, 59, 999);

      // console.log('check event')
      return this.events.some((item) => start <= item.start && end >= item.start)
    },
    getEvents(date) {
      var start = new Date(date);
      start.setHours(0, 0, 0, 0);

      var end = new Date(date);
      end.setHours(23, 59, 59, 999);

      console.log('filter events')
      return this.events.filter((item) => start <= item.start && end >= item.start)
    },
    getEventDurationHours(event) {
      return (event.end - event.start) / (60 * 60 * 1000.)
    }
  },
  mounted() {
    this.updateFirstAndLastDayOfSelection(this.value)
  },
  watch: {
    'value': function () {
      this.updateFirstAndLastDayOfSelection(this.value)
    }
  }
}
</script>

<style>
.theme--light.v-calendar-weekly .v-calendar-weekly__day.v-outside {
  background: none !important;
}

.theme--light.v-calendar-weekly .v-calendar-weekly__head-weekday {
  border: none !important;
  background: none !important;
}

.theme--light.v-calendar-weekly .v-calendar-weekly__day {
  border: none !important;
}

.theme--light.v-calendar-daily .v-calendar-daily__day-interval {
  margin-right: 5px;
  margin-left: 5px;
}

.theme--light.v-calendar-weekly {
  border: none;
}
</style>

<style scoped>
.v-calendar-month--selected-week {
  /*background: var(--v-secondary-base);*/
  background-color: #b0bec51f;

}
</style>