<template>
  <div>
    <v-sheet
        tile
        class="d-flex justify-end pa-5"
    >
      <span class="pa-2"> Projects: </span>
      <div
          v-for="project in projects"
          :key="project.id"
          class="pl-3"
      >
        <v-btn depressed outlined
               :color="project === selectedProject ? 'primary' : ''"
               @click="selectProject(project)">
          {{ project.name }}
        </v-btn>
      </div>

      <span
          v-if="!projects || projects.length === 0"
          class="pa-2 inactive--text"
      >
        Your user is not assign to any project yet
      </span>
    </v-sheet>

    <v-sheet
        tile
        class="d-flex pa-2"
    >
      <v-btn
          icon
          @click="prevWeek()"
      >
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>

      <v-spacer></v-spacer>
      <v-toolbar-title v-if="$refs.calendar">
        {{ $refs.calendar.title }}
      </v-toolbar-title>
      <v-spacer></v-spacer>

      <v-btn
          icon
          @click="nextWeek()"
      >
        <v-icon>mdi-chevron-right</v-icon>
      </v-btn>
    </v-sheet>

    <v-calendar
        class="pa-5"
        ref="calendar"
        v-model="value"
        color="primary"
        type="week"
        first-interval="6"
        interval-count="18"
        interval-height="32"
        interval-minutes="60"
        weekdays="1, 2, 3, 4, 5, 6, 0"
        :events="events"
        :event-color="getEventColor"
        :event-ripple="false"
        @click:event="showEvent"
        @mousedown:event="startDrag"
        @mousedown:time="startTime"
        @mousemove:time="mouseMove"
        @mouseup:time="endDrag"
        @mouseleave.native="cancelDrag"
    >
      <template v-slot:event="{ event, timed, eventSummary }">
        <div
            class="v-event-draggable"
            v-html="eventSummary()"
        ></div>
        <div
            v-if="timed"
            class="v-event-drag-bottom"
            @mousedown.stop="extendBottom(event)"
        ></div>
      </template>

      <template v-slot:day-label-header="{ day, present }">
        <span v-bind:class="present ? 'primary--text' : ''">
          {{ day }}
        </span>
      </template>
    </v-calendar>

    <event-details-dialog v-model="detailsModalOpen"
                          :selected-element="selectedElement" :selected-event="selectedEvent"
                          :projects="projects" :tags="tags"
                          @refreshEvents="$emit('refreshEvents', null)"
    />

    <create-event-dialog v-model="createEventModalOpen"
                         :event="selectedEvent"
                         @refreshEvents="$emit('refreshEvents', null)"
    />
  </div>
</template>

<script>
import EventService from '@/service/EventService'
import EventDetailsDialog from '@/components/workingtime/EventDetailsDialog'
import CreateEventDialog from '@/components/workingtime/CreateEventDialog'
import dayjs from 'dayjs'

const VALUE_FORMAT = 'YYYY-MM-DD'

export default {
  components: {
    EventDetailsDialog,
    CreateEventDialog,
  },

  props: {
    value: {
      type: String,
      required: true,
    },
    selectedTag: {
      type: Object,
      required: false,
    },
    events: {
      type: Array,
      required: true,
    },
    tags: {
      type: Array,
      required: true,
    },
    projects: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    dragEvent: null,
    extendEvent: null,
    createEvent: null,
    selectedProject: null,
    selectedEvent: {},
    selectedElement: null,
    detailsModalOpen: false,
    createEventModalOpen: false,
    keyUpEscapeListener: null,
  }),
  methods: {
    nextWeek() {
      this.$emit('changeCalendarValue', dayjs(this.value, VALUE_FORMAT).add(7, 'days').format(VALUE_FORMAT))
    },
    prevWeek() {
      this.$emit('changeCalendarValue', dayjs(this.value, VALUE_FORMAT).subtract(7, 'days').format(VALUE_FORMAT))
    },
    showEvent({nativeEvent, event}) {

      const open = () => { //FIXME refactor! Something wrong with this code
        this.selectedEvent = event
        this.selectedElement = nativeEvent.target
        setTimeout(() => {
          this.detailsModalOpen = true
        }, 10)
      }

      if (this.detailsModalOpen) {
        this.detailsModalOpen = false
        setTimeout(open, 10)
      } else {
        open()
      }

      nativeEvent.stopPropagation()
    },
    startDrag({event, timed}) {
      console.log('startDrag')
      if (event && timed) {
        this.dragEvent = event
      }
    },
    startTime(tms) {
      console.log('startTime')
      const mouse = this.toTime(tms)

      if (this.detailsModalOpen) {
        console.log('event selected')
        this.detailsModalOpen = false
        return
      }

      if (!this.selectedTag) {
        console.log('tag not selected')
        return
      }

      if (!this.dragEvent) {
        let start = this.roundTime(mouse)

        this.createEvent = {
          name: this.selectedTag.name,
          color: this.selectedTag.color,
          start: start,
          end: start,
          timed: true,
          tag: this.selectedTag,
          project: this.selectedProject,
        }

        this.$emit('addEventToList', this.createEvent)
      }
    },
    extendBottom(event) {
      console.log('extendBottom')
      this.extendEvent = event
    },
    mouseMove(tms) {
      const mouse = this.toTime(tms)

      if (this.extendEvent) {
        const mouseRounded = this.roundTime(mouse, false)

        const min = Math.min(mouseRounded, this.extendEvent.start)
        const max = Math.max(mouseRounded, this.extendEvent.start)

        this.extendEvent.start = min
        this.extendEvent.end = max
      } else if (this.createEvent /*&& this.createStart !== null*/) {
        const mouseRounded = this.roundTime(mouse, false)

        const min = Math.min(mouseRounded, this.createEvent.start)
        const max = Math.max(mouseRounded, this.createEvent.start)

        this.createEvent.start = min
        this.createEvent.end = max
      }
    },
    endDrag() {
      console.log('endDrag')
      if (this.dragEvent || this.createEvent || this.extendEvent) {
        if (this.createEvent) {
          console.log('endDrag - createEvent')
          this.selectedEvent = this.createEvent
          this.createEventModalOpen = true
        }
        if (this.extendEvent) {
          this.updateEvent(this.extendEvent)
        }

        this.dragEvent = null
        this.createEvent = null
        this.extendEvent = null
      }
    },
    cancelDrag() {
      console.log('cancel drug')
      if (this.dragEvent || this.extendEvent || this.createEvent) {
        this.$emit('refreshEvents', null)
      }

      this.createEvent = null
      this.extendEvent = null
      this.dragEvent = null
    },
    roundTime(time, down = true) {
      const roundTo = 15 // minutes
      const roundDownTime = roundTo * 60 * 1000

      return down
          ? time - time % roundDownTime
          : time + (roundDownTime - (time % roundDownTime))
    },
    toTime(tms) {
      return new Date(tms.year, tms.month - 1, tms.day, tms.hour, tms.minute).getTime()
    },
    getEventColor(event) {
      const rgb = parseInt(event.color.substring(1), 16)
      const r = (rgb >> 16) & 0xFF
      const g = (rgb >> 8) & 0xFF
      const b = (rgb >> 0) & 0xFF

      if (event === this.dragEvent) {
        return `rgba(${r}, ${g}, ${b}, 0.7)`
      }

      if (event === this.createEvent) {
        return `rgba(${r}, ${g}, ${b}, 0.7)`
      }

      if (this.selectedProject && event.project !== this.selectedProject) {
        return `rgba(${r}, ${g}, ${b}, 0.4)`
      }

      return event.color
    },
    selectProject(project) {
      if (project === this.selectedProject) {
        this.selectedProject = null
      } else {
        this.selectedProject = project
      }
    },
    saveEvent(event) {
      EventService.create(event)
          .then(() => this.$emit('refreshEvents', null))
    },
    updateEvent(event) {
      EventService.update(event)
          .then(() => this.$emit('refreshEvents', null))
    },
  },
  mounted() {
    this.keyUpEscapeListener = (evt) => {
      if (evt.key === 'Escape') {
        this.detailsModalOpen = false
      }
    }

    document.addEventListener('keyup', this.keyUpEscapeListener)
  },
  beforeDestroy() {
    document.removeEventListener('keyup', this.keyUpEscapeListener)
  },
}
</script>

<style>
.v-calendar-daily__scroll-area {
  overflow-y: hidden !important;
}

.v-calendar-daily__head {
  margin-right: 0 !important;
}

.theme--light.v-calendar-daily {
  border: none !important;
}

.theme--light.v-calendar-daily .v-calendar-daily_head-day {
  border: none !important;
}

.theme--light.v-calendar-daily .v-calendar-daily__intervals-head {
  border: none !important;
}

.theme--light.v-calendar-daily .v-calendar-daily__intervals-head::after {
  background: none !important;
}

.theme--light.v-calendar-daily .v-calendar-daily__day {
  border-right: none !important;
}

</style>

<style scoped lang="scss">
.v-event-draggable {
  padding-left: 6px;
}

.v-event-timed {
  user-select: none;
  -webkit-user-select: none;
}

.v-event-drag-bottom {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 4px;
  height: 4px;
  cursor: ns-resize;

  &::after {
    display: none;
    position: absolute;
    left: 50%;
    height: 4px;
    border-top: 1px solid white;
    border-bottom: 1px solid white;
    width: 16px;
    margin-left: -8px;
    opacity: 0.8;
    content: '';
  }

  &:hover::after {
    display: block;
  }

  body {
    -webkit-touch-callout: none;
    -webkit-user-select: none;
    -khtml-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
  }
}
</style>