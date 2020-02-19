Vue.config.devtools = true;

var store = {
  state : {
    token: '',
    instagramAccounts: [],
    jobs: [],
    settings: {},
    currentUser: null
  },

  setToken (value) { this.state.token = value },
  setInstgramAccounts (value) { this.state.instagramAccounts = value },
  setJobs (value) { this.state.jobs = value },
  setSettings (value) { this.state.settings = value },
  setCurrentUser (value) { this.state.currentUser = value }
}

var loginContainer = new Vue({
  el: '#login-container',
  data: {
    privateState: {
      isLoading: false,
      isFullPage: true,
      username: '',
      password: '',
      loggedIn: false,
      error: null
    }
  },
  created: async () => {
    const token = await window.localStorage.getItem('token');
    
    if (!!token) {
      store.setToken(token);
      loginContainer.getMe();
    }
  },
  methods: {
    getMe: async () => {
      try {  
        var meResponse = await axios.get('{{ .Host }}/api/me', {
          headers: {
            Authorization: `Bearer ${store.state.token}`
          }
        });
        
        var me = meResponse.data;

        store.setCurrentUser(me);
        loginContainer.privateState.loggedIn = true;
        settings.sharedState.enabled = true;
        settings.refreshTickets();
        settings.setSettings(me.Settings);
        settings.setInstgramAccounts(me.InstagramAccounts);
      } catch (err) {
        loginContainer.privateState.error = err;
      }
    },
    logout: async () => {
      await window.localStorage.removeItem('token');
      loginContainer.privateState.loggedIn = false;
      settings.sharedState.enabled = false;
    },
    authenticate: async () => {
      loginContainer.privateState.isLoading = true;

      try {
        var tokenResponse = await axios.post('{{ .Host }}/auth', {
            username: loginContainer.privateState.username,
            password: loginContainer.privateState.password
        });


        store.setToken(tokenResponse.data.AccessToken);
        await window.localStorage.setItem(`token`, tokenResponse.data.AccessToken);

        await loginContainer.getMe();
      } catch (err) {
        loginContainer.privateState.error = err;
      }

      loginContainer.privateState.isLoading = false;
    }
  }
});

var settings = new Vue({
  el: '#settings',
  methods: {
    setSettings: (newSettings) => {
      settings.privateState.settings = newSettings;
    },
    setInstgramAccounts: (accounts) => {
      settings.privateState.instagramAccounts = accounts;
      settings.privateState.selectedInstagramAccount = accounts[0];
    },
    refreshTickets: async () => {
      const ticketsResponse = await axios.get('{{ .Host }}/api/tickets', {
        headers: {
          Authorization: `Bearer ${store.state.token}`
        }
      });

      const tickets = ticketsResponse.data;

      settings.privateState.runningJobs = tickets.filter((t) => !t.Done);
      settings.privateState.finishedJobs = tickets.filter((t) => t.Done);
    },
    displayLogs: async (id) => {
      const logsResponse = await axios.get(`{{ .Host }}/api/tickets/${id}/logs`, {
        headers: {
          Authorization: `Bearer ${store.state.token}`
        }
      });

      const logs = logsResponse.data;

      logs.ErrLogs = logs.ErrLogs.map((s) => s.trim(' ')).filter(s => !s.includes('Description:'));
      Buefy.ModalProgrammatic.open({
        parent: settings,
        component: LogsModal,
        'custom-class': 'container',
        props: logs
      })
    }, 
    saveSettings: async () => {
      try {
        await axios.post('{{ .Host }}/api/settings', { Settings: settings.privateState.settings }, {
          headers: {
            Authorization: `Bearer ${store.state.token}`
          }
        });
        Buefy.ToastProgrammatic.open({
          duration: 5000,
          message: 'Settings saved',
          type: 'is-success'
        })
      } catch (err) {
        Buefy.ToastProgrammatic.open({
          duration: 5000,
          message: `Could not save: ${err.response.data.error}`,
          type: 'is-danger'
        });
      }
    },
    runJob: async () => {
      try {
        runJobResponse = await axios.post('{{ .Host }}/api/jobs', {
          Label: "job",
          IGPassword: settings.privateState.igPassword,
          Settings: {
            Account: settings.privateState.selectedInstagramAccount,
            Settings: settings.privateState.settings
          }
        }, {
          headers: {
            Authorization: `Bearer ${store.state.token}`
          }
        });
        await settings.refreshTickets();
      } catch (err) {
        Buefy.ToastProgrammatic.open({
          duration: 5000,
          message: `Could not run the bot: ${err.response.data.error}`,
          type: 'is-danger'
        });
      }
    }
  },
  data: {
    sharedState: {
      enabled: false
    },
    privateState: {
      settings: {
        Hashtags: [],
        Comments: [],
        TotalLikes: 0,
        PotencyMode: "negative",
        MaxFollowers: 0,
        MinFollowers: 0,
        MaxFollowing: 0,
        MinFollowing: 0
      },
      instagramAccounts: [],
      selectedInstagramAccount: null,
      igPassword: null,
      runningJobs: [],
      finishedJobs: []
    }
  }
})

