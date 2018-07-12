new Vue({
    el: '#app3',
    data: {
        info: null
    },
    created() {
        fetch('http://localhost:9000/') //fetching JSON data from page
            .then(response=> (this.info = response))
    }
})