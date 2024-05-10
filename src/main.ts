import './style/main.scss';
import * as bootstrap from 'bootstrap';
import Swal from 'sweetalert2';
import ckeditor from './ckeditor';
import htmx from './htmx';
// import * as dropzone from 'dropzone';
import Alpine from 'alpinejs';

declare global {
  interface Window {
    htmx: typeof htmx;
    ckeditor: typeof ckeditor;
    bootstrap: typeof bootstrap;
    Swal: typeof Swal;
    Alpine: typeof Alpine;
  }
}

window.htmx = htmx;
window.ckeditor = ckeditor;
window.bootstrap = bootstrap;
window.Swal = Swal;
window.Alpine = Alpine;

Alpine.data('date', () => ({
  formatDate(unix: any) {
    const unixTimestampMilliseconds = parseInt(unix) * 1000;
    const date = new Date(unixTimestampMilliseconds);
    const localeOptions: any = {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    };
    return date.toLocaleDateString('id-ID', localeOptions);
  },
}));

Alpine.start();
