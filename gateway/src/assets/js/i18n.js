var i18n = {
    "locale": "it",
    setLocale(locale) { this.locale = locale },
    get(key) { return this[this.locale][key] },
    add(locale, dictionary) { this[locale] = dictionary; },
    init(dicts) { for (locale in dicts) { i18n.add(locale, dicts[locale]) }; },
    apply() { for (key in this[this.locale]) { try { 
      var subkeys = key.split(".");
      if (subkeys.length == 1)
        document.getElementById(key).innerHTML = this.get(key); // <div id="<key>"></div>
      else if (subkeys[0] == "meta")
        document[subkeys[1]] = this.get(key); // meta
      else
        document.getElementsByTagName(subkeys[1])[0].innerHTML = this.get(key); // HTML tagName: H1, ...
    } catch (e) { } } }
  };