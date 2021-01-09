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
                'inactive--text': !isCurrentMonth(month),
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
          <template>
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
import dayjs from 'dayjs'

const VALUE_FORMAT = 'YYYY-MM-DD'

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
    days: {
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
      let selectedDate = new Date(date)
      return selectedDate >= this.selectedWeekFirstDay && selectedDate <= this.selectedWeekLastDay
    },
    updateFirstAndLastDayOfSelection(date) {
      let selectedDate = dayjs(date, VALUE_FORMAT).utc(true)
      this.selectedMonth = selectedDate.month()

      this.selectedWeekFirstDay = selectedDate.day(1).toDate()
      this.selectedWeekLastDay = selectedDate.day() === 0 ? selectedDate : selectedDate.day(7).toDate()
    },
    isCurrentMonth(month) {
      return this.selectedMonth === (month - 1)
    },
    getEvents(date) {
      const selectedDate = dayjs(date, VALUE_FORMAT).utc(true)

      const day = this.days.find((day) => selectedDate.isSame(day.date))
      return day ? day.events : []
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
  background-color: #b0bec51f;
}
</style>