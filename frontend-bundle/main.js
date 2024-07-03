import './scss/main.scss';
import 'htmx.org';
import './pkg/htmx-ext'
import Swal from 'sweetalert2';
import ckeditor from './pkg/ckeditor';
import * as bootstrap from 'bootstrap';
import Alpine from 'alpinejs';
import { formatDate } from './formatdate';
import focus from '@alpinejs/focus';

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