<template>
  <v-menu
      v-model="dateRangeMenu"
      :close-on-content-click="false"
      transition="scale-transition"
      offset-y
  >
    <template v-slot:activator="{ on, attrs }">
      <v-text-field
          v-model="dateRangeText"
          prepend-icon="mdi-calendar"
          v-bind="attrs"
          @blur="date = parseDate(value)"
          v-on="on"
      ></v-text-field>
    </template>
    <v-date-picker
        :value="value"
        @input="emitChange"
        @change="$emit('change', '')"
        no-title
        range
    ></v-date-picker>
  </v-menu>
</template>

<script>
import dayjs from "dayjs";

export default {
  name: "DateRangePicker",

  props: {
    value: {
      type: Array,
      required: false,
    },
  },
  data: () => ({
    dateRangeMenu: false,
  }),
  methods: {
    parseDate(date) {
      if (!date) return null

      return dayjs(date, 'M/D/YYYY').format('YYYY-MM-DD')
    },
    emitChange(value) {
      this.$emit('input', value)
    },
  },
  computed: {
    dateRangeText() {
      if (!this.value || this.value.length === 0) {
        return ''
      }

      return this.value.join(' - ')
    },
  },
}
</script>