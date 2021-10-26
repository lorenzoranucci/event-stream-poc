<?php

declare(strict_types=1);

namespace App\MerchantReview\Infrastructure;

use App\MerchantReview\Domain\MerchantReviewUUID;
use App\MerchantReview\Domain\MerchantReviewUUIDGenerator;

final class DummyMerchantReviewUUIDGenerator implements MerchantReviewUUIDGenerator
{
    public function generate(): MerchantReviewUUID
    {
        return MerchantReviewUUID::createFromString(uniqid());
    }
}
