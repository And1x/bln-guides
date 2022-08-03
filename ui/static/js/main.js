function jsToggleDayNight() {
    if (localStorage == null) {
        localStorage.setItem('darkmode', 'true');
    }

    const theme = document.getElementById('cssstyle');

    if (localStorage.getItem('darkmode') == 'true' ){
        theme.setAttribute('href', '/static/css/day.css');
        localStorage.SetItem('darkmode', 'false');
        console.log("was here");
    } else {
        theme.setAttribute('href', '/static/css/main.css');
        localStorage.SetItem('darkmode', 'true');
        console.log("heree>?");
    }
 }