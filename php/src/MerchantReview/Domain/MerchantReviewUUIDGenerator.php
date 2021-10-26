<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

interface MerchantReviewUUIDGenerator
{
    public function generate(): MerchantReviewUUID;
}
