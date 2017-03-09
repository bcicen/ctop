/* global GMaps: true */

$(document).ready(function () {
  map()
})

function map () {
  if ($('#map').length) {
    var lat = $('#gmap-lat').val()
    var lng = $('#gmap-lng').val()
    var image = '/img/marker.png'

    var styles =
      [
        {
          'featureType': 'landscape', 'stylers': [{'saturation': -100}, {'lightness': 65}, {'visibility': 'on'}]
        }, {
          'featureType': 'poi', 'stylers': [{'saturation': -100}, {'lightness': 51}, {'visibility': 'simplified'}]
        }, {
          'featureType': 'road.highway', 'stylers': [{'saturation': -100}, {'visibility': 'simplified'}]
        }, {
          'featureType': 'road.arterial', 'stylers': [{'saturation': -100}, {'lightness': 30}, {'visibility': 'on'}]
        }, {
          'featureType': 'road.local', 'stylers': [{'saturation': -100}, {'lightness': 40}, {'visibility': 'on'}]
        }, {
          'featureType': 'transit', 'stylers': [{'saturation': -100}, {'visibility': 'simplified'}]
        }, {
          'featureType': 'administrative.province', 'stylers': [{'visibility': 'off'}]
        }, {
          'featureType': 'water', 'elementType': 'labels', 'stylers': [{'visibility': 'on'}, {'lightness': -25}, {'saturation': -100}]
        }, {
          'featureType': 'water', 'elementType': 'geometry', 'stylers': [{'hue': '#ffff00'}, {'lightness': -25}, {'saturation': -97}]
        }
      ]

    var map = new GMaps({
      el: '#map',
      lat: lat,
      lng: lng,
      zoomControl: true,
      zoomControlOpt: {
        style: 'SMALL',
        position: 'TOP_LEFT'
      },
      panControl: false,
      streetViewControl: false,
      mapTypeControl: false,
      overviewMapControl: false,
      scrollwheel: false,
      draggable: false,
      styles: styles
    })

    map.addMarker({
      lat: lat,
      lng: lng,
      icon: image /* ,
      title: '',
      infoWindow: {
      content: '<p>HTML Content</p>'
      } */
    })
  }
}
