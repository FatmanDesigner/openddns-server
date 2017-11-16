import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return this.get('store').findAll('domain-entry');
  },
  actions: {
    error(error, transition) { // eslint-disable-line
      const httpError = error.errors && error.errors.filter(item => item.status==='401');
      if (httpError) {
        this.replaceWith('login')
      } else {
        return true;
      }
    }
  }
});
