<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

interface MerchantReviewRepository
{
    public function add(MerchantReview $merchantReview): void;

    public function getOneByUUID(MerchantReviewUUID $merchantReviewUUID): MerchantReview;
}
