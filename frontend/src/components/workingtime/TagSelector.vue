<template>
  <div class="pa-5">
    <v-row
        v-for="tag in tags"
        :key="tag.id"
        align="center"
        justify="start"
        class="ma-1"
    >
      <v-btn depressed outlined :color="tag === selectedTag ? 'primary' : ''" @click="selectTag(tag)">
        Add
      </v-btn>
      <span class="pl-5" v-bind:class="tag === selectedTag ? 'primary--text' : ''">
      {{ tag.name }}
      </span>
      <v-spacer></v-spacer>
      <v-icon v-bind:color="tag.color">mdi-format-color-highlight</v-icon>
    </v-row>
  </div>
</template>

<script>
export default {
  components: {},
  props: {
    tags: {
      type: Array,
      required: true,
    },
    selectedTag: {
      type: Object,
      required: false,
    },
  },
  data: () => ({}),
  methods: {
    selectTag(tag) {
      if (tag === this.selectedTag) {
        this.$emit('changeSelectedTag', null)
      } else {
        this.$emit('changeSelectedTag', tag)
      }
    },
  }
}
</script>