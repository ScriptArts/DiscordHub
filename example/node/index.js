const NATS = require('nats')
const nats = NATS.connect()

nats.subscribe('SAMPLE_SUBJECT', function (request, replyTo) {
    console.log(request)

    const data = JSON.stringify({
        type: 0,
        message: 'foobar!!'
    })

    nats.publish(replyTo, data)
})
