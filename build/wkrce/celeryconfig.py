import os

broker_url = os.environ.get(
    'BROKER_URI',
    'amqp://surge:paragon@queue:5672/celery'
)

task_serializer = 'json'
accept_content = ['json']

imports = ('tasks', )

task_routes = {
    'fibonacci': os.environ.get('WORKER_QUEUE', 'celery'),
}

task_default_queue = 'celery'
task_default_exchange = 'default'
task_default_routing_key = 'default'
celery_task_protocol = 1
