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
        :first-interval="firstInterval"
        :interval-count="intevalCount"
        :interval-height="calcIntervalHeight()"
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
import cloneDeep from 'clone-deep'

const VALUE_FORMAT = 'YYYY-MM-DD'
const MIN_EVENT_HEIGHT = 24
const MAX_EVENT_HEIGHT = 48
const DEFAULT_EVENT_HEIGHT = 32
const DEFAULT_WEAK_HOURS = 10

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
    selectedProject: {
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
    selectedEvent: {},
    selectedElement: null,
    detailsModalOpen: false,
    createEventModalOpen: false,
    keyUpEscapeListener: null,

    firstInterval: 6,
    intevalCount: 18,
  }),
  methods: {
    calcIntervalHeight() {
      if (localStorage.minTimeRange || localStorage.maxTimeRange) {
        const min = parseInt(localStorage.minTimeRange)
        const max = parseInt(localStorage.maxTimeRange)

        const calculatedHeight = DEFAULT_EVENT_HEIGHT * DEFAULT_WEAK_HOURS / (max - min);

        if (calculatedHeight >= MIN_EVENT_HEIGHT && calculatedHeight <= MAX_EVENT_HEIGHT) {
          return calculatedHeight
        } else {
          if (calculatedHeight < MIN_EVENT_HEIGHT) {
            return MIN_EVENT_HEIGHT
          }
          if (calculatedHeight > MAX_EVENT_HEIGHT) {
            return MAX_EVENT_HEIGHT
          }
        }
      }

      return DEFAULT_EVENT_HEIGHT
    },
    nextWeek() {
      this.$emit('changeCalendarValue', dayjs(this.value, VALUE_FORMAT).add(7, 'days').format(VALUE_FORMAT))
    },
    prevWeek() {
      this.$emit('changeCalendarValue', dayjs(this.value, VALUE_FORMAT).subtract(7, 'days').format(VALUE_FORMAT))
    },
    showEvent({nativeEvent, event}) {

      const open = () => { //FIXME refactor! Something wrong with this code
        this.selectedEvent = cloneDeep(event)
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
      if (event && timed) {
        this.dragEvent = event
      }
    },
    startTime(tms) {
      const mouse = this.toTime(tms)

      if (this.detailsModalOpen) {
        this.detailsModalOpen = false
        return
      }

      if (!this.selectedTag) {
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
      if (this.dragEvent || this.createEvent || this.extendEvent) {
        if (this.createEvent) {
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
        this.$emit('changeSelectedProject', null)
      } else {
        this.$emit('changeSelectedProject', project)
      }
    },
    saveEvent(event) {
      EventService.create(event)
          .then(() => this.$emit('refreshEvents', null))
    },
    updateEvent(event) {
      let e = cloneDeep(event);

      let start = dayjs(e.start)
      let end = dayjs(e.end)

      if (!end.isAfter(start)) {
        end = start.add(15, 'minute')
      }

      if (start.day() !== end.day() || start.month() !== end.month() || start.day() !== start.day()) {
        end = start.endOf('day')
      }

      e.start = start.toISOString()
      e.end = end.toISOString()

      EventService.update(e)
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

    if (localStorage.minTimeRange && localStorage.maxTimeRange) {
      const min = parseInt(localStorage.minTimeRange)
      const max = parseInt(localStorage.maxTimeRange)

      this.firstInterval = min == 0 ? 0 : min - 1
      
      const count = max - min + 2
      this.intevalCount = count > 24 ? 24 : count
    }
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