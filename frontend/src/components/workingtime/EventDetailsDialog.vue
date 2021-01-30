<template>
  <v-menu
      :value="value"
      @input="emitChange"
      :close-on-content-click="false"
      :activator="selectedElement"
      offset-x
  >
    <v-card
        color="grey lighten-4"
        min-width="350px"
        flat
    >
      <v-form
          v-model="selectedEventValid"
      >
        <v-toolbar
            :color="selectedEvent.color"
            dark
        >
          <v-text-field
              label="name"
              hide-details="auto"
              v-model="selectedEvent.name"
              @change="updateEvent(selectedEvent)"
              required
          ></v-text-field>
        </v-toolbar>
        <v-card-text>
          <v-text-field
              label="from"
              hide-details="auto"
              :rules="timeRules"
              :value="getTime(selectedEvent.start)"
              @change="updateStartTime"
              required
          ></v-text-field>
          <v-text-field
              label="to"
              hide-details="auto"
              :rules="timeRules"
              :value="getTime(selectedEvent.end)"
              @change="updateEndTime"
              required
          ></v-text-field>
          <v-select
              label="Tag"
              class="mt-7"
              :items="tags"
              v-model="selectedEvent.tag"
              v-on:change="updateEvent(selectedEvent)"
              item-text="name"
              return-object
              outlined
              clearable
          />
          <v-select
              label="Project"
              :items="projects"
              v-model="selectedEvent.project"
              v-on:change="updateEvent(selectedEvent)"
              item-text="name"
              return-object
              outlined
              clearable
          />
        </v-card-text>
      </v-form>

      <div class="d-flex justify-center">
        <div class="pl-3 pr-3 global-form-error error--text">
          {{ globalFormError }}
        </div>
      </div>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
            icon
            @click="deleteEvent(selectedEvent)"
        >
          <v-icon>mdi-trash-can-outline</v-icon>
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script>

import EventService from "@/service/EventService"

export default {
  name: "EventDetailsDialog",

  props: {
    value: {
      type: Boolean,
      required: false,
    },
    selectedElement: {
      type: HTMLDivElement,
      required: false,
    },
    selectedEvent: {
      type: Object,
      required: false,
    },
    projects: {
      type: Array,
      required: true,
    },
    tags: {
      type: Array,
      required: true,
    },
  },
  data: () => ({
    selectedEventValid: true,
    globalFormError: '',
    timeRules: [
      value => !!value || 'Required.',
      value => value.split(':').length <= 2 && value.split(':')[0] < 24 || 'Wrong hours',
      value => (value.split(':').length < 2 || (value.split(':').length === 2 && value.split(':')[1] < 60)) || 'Wrong minutes',
    ],
  }),
  methods: {
    getTime(date) {
      let d = new Date(date);
      let minutes = d.getMinutes();
      let hours = d.getHours();

      return (hours >= 10 ? hours : `0${hours}`) + ':' + (minutes >= 10 ? minutes : `0${minutes}`)
    },
    getDate(date, time) {
      let d = new Date(date)
      let timeElements = time.split(':');

      d.setHours(timeElements[0])
      if (timeElements.length === 2) {
        d.setMinutes(timeElements[1])
      } else {
        d.setMinutes(0)
      }
      return d
    },
    updateStartTime(e) {
      if (this.selectedEventValid) {
        this.selectedEvent.start = this.getDate(this.selectedEvent.start, e)
        this.updateEvent(this.selectedEvent)
      }
    },
    updateEndTime(e) {
      if (this.selectedEventValid) {
        this.selectedEvent.end = this.getDate(this.selectedEvent.end, e)
        this.updateEvent(this.selectedEvent)
      }
    },
    updateEvent(event) {
      EventService.update(event)
          .then(() => {
            this.$emit('refreshEvents', null)
          })
          .catch(error => {
            if (error && error.response && (error.response.status === 400 || error.response.status === 404)) {
              this.globalFormError = error.response.data.error
            }
          })
    },
    deleteEvent(event) {
      this.emitChange(false)

      EventService.delete(event)
          .then(() => this.$emit('refreshEvents', null))
    },
    emitChange(value) {
      this.$emit('input', value)
    },
  },
  watch: {
    'value': function () {
      this.globalFormError = null
    },
  },
}
</script>

<style scoped>
.global-form-error {
  width: 350px;
  text-align: center;
  white-space: pre-line;
}
</style>