<template>
  <div class="pa-5">
    <div
      v-for="projectTags in groupByProject(tags)"
      :key="projectTags[0].project"
    >
      <span v-if="findProjectNameById(projectTags[0].project)" class="inactive--text"> 
        {{ findProjectNameById(projectTags[0].project).name }} 
      </span>
      
      <v-row
        v-for="tag in projectTags"
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
    projects: {
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
    groupByProject(tags) {
      console.log('call groupBys');
      return tags.reduce((r, a) => {
        r[a.project] = [...r[a.project] || [], a];
        return r;
      }, {})
    },
    findProjectNameById(id) {
      if (this.projects) {
        return this.projects.find(p => p.id === id)
      }
    },
  }
}
</script>