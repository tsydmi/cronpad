<template>
  <div class="pr-5 pl-5 pt-3 pb-3">
    <v-row
        class="pb-1"
    >
      <v-dialog
          v-model="createTagDialog"
          max-width="400px"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
              icon
              v-bind="attrs"
              v-on="on"
          >
            <v-icon>mdi-plus</v-icon>
          </v-btn>
        </template>
        <CreateTag :dialog="createTagDialog" :tags="tags" @closeDialog="createTagDialog = false"
                   @refreshTags="refreshTags"/>
      </v-dialog>
      <v-spacer></v-spacer>
      <v-dialog
          v-model="configureTagDialog"
          max-width="400px"
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn
              icon
              v-bind="attrs"
              v-on="on"
          >
            <v-icon>mdi-cog-outline</v-icon>
          </v-btn>
        </template>
        <ConfigureTags :dialog="configureTagDialog" :tags="tags" @closeDialog="configureTagDialog = false"
                       @refreshTags="refreshTags"/>
      </v-dialog>
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
      <!--      <v-icon v-bind:color="tag.color">mdi-tag</v-icon>-->
    </v-row>

  </div>
</template>

<script>
import CreateTag from "@/components/tag/CreateTag";
import ConfigureTags from "@/components/tag/ConfigureTags";

export default {
  components: {CreateTag, ConfigureTags},
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
    createTagDialog: false,
    configureTagDialog: false,
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