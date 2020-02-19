const LogsModal = {
  props: ['Logs', 'ErrLogs'],
  template: `
    <pre class="section white">
      <p v-for="line in ErrLogs">{{ line }}</p>
    </pre>
  `
}
