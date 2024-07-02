function formatDate(unixTimestamp) {
  const date = new Date(unixTimestamp * 1000);
  const day = date.toLocaleDateString('id-ID', { weekday: 'long' });
  const dayNumber = date.getDate();
  const monthNames = [
    'Januari',
    'Februari',
    'Maret',
    'April',
    'Mei',
    'Juni',
    'Juli',
    'Agustus',
    'September',
    'Oktober',
    'November',
    'Desember',
  ];
  const month = monthNames[date.getMonth()];
  const year = date.getFullYear();
  const hour = ('0' + date.getHours()).slice(-2);
  const minute = ('0' + date.getMinutes()).slice(-2);
  const timeZone = getTimeZoneAbbreviation();

  return `${day}, ${dayNumber} ${month} ${year} | ${hour}:${minute} ${timeZone}`;
}

function getTimeZoneAbbreviation() {
  const offsetMinutes = new Date().getTimezoneOffset();
  const timeZones = {
    '-480': 'WITA', // UTC+8 (Central Indonesian Time)
    '-420': 'WIB', // UTC+7 (Western Indonesian Time)
    '-540': 'WIT', // UTC+9 (Eastern Indonesian Time)
  };

  return timeZones[offsetMinutes] || `UTC${-offsetMinutes / 60}`;
}

export { formatDate };
