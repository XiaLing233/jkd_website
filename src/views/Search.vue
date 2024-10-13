<template>
<main>
    <input
        type="text"
        v-model="searchQuery"
        placeholder="Search for keywords"
    />
    <table>
        <thead>
        <tr>
            <th @click="sortTable('name')">Name</th>
            <th @click="sortTable('age')">Age</th>
            <th @click="sortTable('city')">City</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="item in filteredAndSortedItems" :key="item.id">
            <td>{{ item.name }}</td>
            <td>{{ item.age }}</td>
            <td>{{ item.city }}</td>
        </tr>
        </tbody>
    </table>
    </main>
</template>

<script>
export default {
  data() {
    return {
      searchQuery: "",
      sortKey: "",
      sortAsc: true,
      items: [
        { id: 1, name: "Alice", age: 25, city: "New York" },
        { id: 2, name: "Bob", age: 30, city: "San Francisco" },
        { id: 3, name: "Charlie", age: 28, city: "Los Angeles" },
        // Add more items as needed
      ],
    };
  },
  computed: {
    filteredAndSortedItems() {
      // Filter items by search query
      let filteredItems = this.items.filter((item) => {
        return Object.values(item).some((value) =>
          String(value).toLowerCase().includes(this.searchQuery.toLowerCase())
        );
      });

      // Sort items
      if (this.sortKey) {
        filteredItems.sort((a, b) => {
          let result = a[this.sortKey] > b[this.sortKey] ? 1 : -1;
          return this.sortAsc ? result : -result;
        });
      }
      return filteredItems;
    },
  },
  methods: {
    sortTable(key) {
      if (this.sortKey === key) {
        this.sortAsc = !this.sortAsc; // Toggle sorting order
      } else {
        this.sortKey = key;
        this.sortAsc = true; // Default to ascending order
      }
    },
  },
};
</script>

<style>

</style>