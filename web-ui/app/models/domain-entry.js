import DS from 'ember-data';

export default DS.Model.extend({
  domainName: DS.attr('string'),
  ip: DS.attr('string'),
  updatedAt: DS.attr('date')
});
