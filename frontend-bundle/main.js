import './scss/main.scss';
import 'htmx.org';
import Swal from 'sweetalert2';
import ckeditor from './pkg/ckeditor';
import * as bootstrap from 'bootstrap';
import Alpine from 'alpinejs';
import { formatDate } from './formatdate';
import focus from '@alpinejs/focus';

htmx.defineExtension('path-params', {
  onEvent: function (name, evt) {
    if (name === 'htmx:configRequest') {
      evt.detail.path = evt.detail.path.replace(
        /{([^}]+)}/g,
        function (_, param) {
          var val = evt.detail.parameters[param];
          delete evt.detail.parameters[param];
          return val === undefined
            ? '{' + param + '}'
            : encodeURIComponent(val);
        }
      );
    }
  },
});

window.htmx = htmx;
window.ckeditor = ckeditor;
window.bootstrap = bootstrap;
window.Swal = Swal;
window.Alpine = Alpine;

Alpine.data('date', () => ({
  formatDate,
}));


Alpine.plugin(focus);
Alpine.start();