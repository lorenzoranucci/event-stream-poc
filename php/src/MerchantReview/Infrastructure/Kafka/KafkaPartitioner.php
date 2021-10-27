<?php

declare(strict_types=1);

namespace App\MerchantReview\Infrastructure\Kafka;

class KafkaPartitioner
{
    public function getPartitionID(): int
    {
        return 0;
    }
}
