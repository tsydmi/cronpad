<template>
  <div>
    <span class="v-label v-text-field mr-3">Color</span>
    <v-color-picker
        label="Color"
        hide-canvas
        show-swatches
        :swatches="swatches"
        swatches-max-height="400px"
        :value="value"
        @change="$emit('change', '')"
        @input="emitChange"
    ></v-color-picker>
  </div>
</template>

<script>
export default {
  name: "TagColorPicker",
  props: {
    value: {
      required: true,
    },
  },
  data: () => ({
    swatches: [
      ["#B39DDB", "#9FA8DA", "#90CAF9"],
      ["#81D4FA", "#80DEEA", "#80CBC4"],
      ["#A5D6A7", "#C5E1A5", "#E6EE9C"],
      ["#FFE082", "#FFCC80", "#FFAB91"],
    ],
  }),
  methods: {
    emitChange(value) {
      this.$emit('input', value)
    },
  },
}
</script>