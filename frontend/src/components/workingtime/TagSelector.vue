<template>
  <div class="pr-5 pl-5 pt-3 pb-3">
    <v-row
        class="pb-1"
    >
      <v-btn
          icon
          @click="openCreateTagDialog = true"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>

      <v-spacer></v-spacer>

      <v-btn
          icon
          @click="openConfigureTagDialog = true"
      >
        <v-icon>mdi-cog-outline</v-icon>
      </v-btn>
    </v-row>
    <v-row
        v-for="tag in tags"
        v-bind:key="tag.id"
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

    <create-tag-dialog v-model="openCreateTagDialog"
                       :tags="tags"
                       @refreshTags="refreshTags"/>

    <configure-tags-dialog v-model="openConfigureTagDialog"
                           :tags="tags"
                           @refreshTags="refreshTags"/>

  </div>
</template>

<script>
import CreateTagDialog from "@/components/workingtime/tagselector/CreateTagDialog"
import ConfigureTagsDialog from "@/components/workingtime/tagselector/ConfigureTagsDialog"

export default {
  components: {
    CreateTagDialog,
    ConfigureTagsDialog,
  },
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
  data: () => ({
    openCreateTagDialog: false,
    openConfigureTagDialog: false,
  }),
  methods: {
    selectTag(tag) {
      if (tag === this.selectedTag) {
        this.$emit('changeSelectedTag', null)
      } else {
        this.$emit('changeSelectedTag', tag)
      }
    },
    refreshTags() {
      this.$emit('refreshTags', null)
    },
    configureTags() {

    },
  }
}
</script>