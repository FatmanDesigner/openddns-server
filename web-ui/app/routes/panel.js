import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return [
      {
        domainName: 'Cell',
        ip: 'Cell',
        lastUpdated: 'Cell'
      },
      {
        domainName: 'Cell',
        ip: 'Cell',
        lastUpdated: 'Cell'
      }
    ]
  }
});
