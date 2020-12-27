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
        <MonthCalendar :today="today" :value="value" :events="events"
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
import TagSelector from "@/components/workingtime/TagSelector"
import TagService from '@/service/TagService'
import DayService from '@/service/DayService'
import ProjectService from "@/service/ProjectService"

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
    selectedTag: null,
    selectedProject: null
  }),
  methods: {
    updateValue(value) {
      this.value = value
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
      let date = new Date(this.value)

      var firstDayOfMonth = new Date(Date.UTC(date.getFullYear(), date.getMonth(), 1));
      var lastDayOfMonth = new Date(Date.UTC(date.getFullYear(), date.getMonth() + 1, 0));

      firstDayOfMonth.setDate(firstDayOfMonth.getDate() - this.getRealDayOfWeek(firstDayOfMonth))
      lastDayOfMonth.setDate(lastDayOfMonth.getDate() + (6 - this.getRealDayOfWeek(lastDayOfMonth)))

      DayService.findByDayRange(firstDayOfMonth, lastDayOfMonth)
          .then(response => this.events = this.convertToEvents(response.data));
    },
    convertToEvents(days) {
      let events = []

      days.forEach(day => {
            day.events.forEach(event => {
              let tag = this.getTagById(event.tag);

              events.push({
                id: event.id,
                name: event.name,
                start: new Date(event.start),
                end: new Date(event.end),
                timed: event.timed,
                tag: tag,
                project: this.getProjectById(event.project),
                color: tag ? tag.color : '#7d7d7d' // TODO replace by something global
              })
            })
          }
      )
      return events
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
    let date = new Date()
    this.today = date.toISOString().split('T')[0]
    this.value = date.toISOString().split('T')[0]

    this.refreshTags()
    this.refreshProjects()
    this.refreshEvents()
  }
};
</script>
