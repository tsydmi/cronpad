<template>
  <v-row dense>
    <v-col lg="9" md="8" sm="12">
      <v-card elevation="4">
        <WeekCalendar
            :events="events" :projects="projects" :tags="tags"
            :value="value" :selected-tag="selectedTag"
            @addEventToList="addEventToList"
            @refreshEvents="refreshEvents"
            @changeCalendarValue="updateValue"/>
      </v-card>
    </v-col>
    <v-col lg="3" md="4" sm="12" class="pl-4 pr-4">
      <v-card elevation="4">
        <MonthCalendar :today="today" :value="value" :days="days"
                       @changeCalendarValue="updateValue"/>
        <TagSelector :tags="tags" :selected-tag="selectedTag"
                     @changeSelectedTag="updateSelectedTag"
                     @refreshTags="refreshTags"/>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import WeekCalendar from '@/components/workingtime/WeekCalendar'
import MonthCalendar from '@/components/workingtime/MonthCalendar'
import TagSelector from '@/components/workingtime/TagSelector'
import TagService from '@/service/TagService'
import DayService from '@/service/DayService'
import ProjectService from '@/service/ProjectService'
import dayjs from 'dayjs'

const VALUE_FORMAT = 'YYYY-MM-DD'

export default {
  name: 'WorkingTime',

  components: {
    WeekCalendar,
    MonthCalendar,
    TagSelector
  },

  data: () => ({
    today: null,
    value: null,
    selectedWeekFirstDay: null,
    selectedWeekLastDay: null,
    selectedDay: null,
    tags: [],
    projects: [],
    events: [],
    days: [],
    selectedTag: null,
    selectedProject: null,
  }),
  methods: {
    updateValue(value) {
      const newValueMonth = dayjs(value, VALUE_FORMAT).month()
      const prevValueMonth = dayjs(this.value, VALUE_FORMAT).month()

      this.value = value

      if (newValueMonth !== prevValueMonth) {
        this.refreshEvents()
      }
    },
    updateSelectedTag(tag) {
      this.selectedTag = tag
    },
    addEventToList(event) {
      this.events.push(event)
    },
    refreshTags() {
      TagService.findAll()
          .then(response => this.tags = response.data);
      this.selectedTag = null
    },
    refreshProjects() {
      ProjectService.findCurrentUserProjects()
          .then(response => this.projects = response.data);
    },
    refreshEvents() {
      const date = dayjs(this.value, VALUE_FORMAT)

      const firstDayOfMonth = date.date(1)
      const lastDayOfMonth = date.date(date.daysInMonth())

      const from = firstDayOfMonth.day(1)
      const to = lastDayOfMonth.day() === 0 ? lastDayOfMonth : lastDayOfMonth.day(7)

      DayService.findByDayRange(from, to)
          .then(response => this.convertEvents(response.data))
    },
    convertEvents(days) {
      let convertedEvents = []
      let convertedDays = []

      days.forEach(day => {
            let convertedDay = {
              date: dayjs(day.date).utc(),
              events: [],
            }

            day.events.forEach(event => {
              let tag = this.getTagById(event.tag);

              const convertedEvent = {
                id: event.id,
                name: event.name,
                start: new Date(event.start),
                end: new Date(event.end),
                timed: event.timed,
                tag: tag,
                project: this.getProjectById(event.project),
                color: tag ? tag.color : '#7d7d7d' // TODO replace by something global
              }
              convertedEvents.push(convertedEvent)
              convertedDay.events.push(convertedEvent)
            })
            convertedDays.push(convertedDay)
          }
      )

      this.days = convertedDays
      this.events = convertedEvents
    },
    getTagById(tagId) {
      return this.tags.find(tag => tag.id === tagId);
    },
    getProjectById(projectId) {
      return this.projects.find(project => project.id === projectId)
    },
    getRealDayOfWeek(date) {
      return (date.getDay() + 6) % 7
    }
  },
  created() {
    let date = dayjs()
    this.today = date.format(VALUE_FORMAT)
    this.value = date.format(VALUE_FORMAT)

    this.refreshTags()
    this.refreshProjects()
    this.refreshEvents()
  }
};
</script>
