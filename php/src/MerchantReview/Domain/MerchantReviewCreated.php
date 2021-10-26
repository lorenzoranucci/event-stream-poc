<?php

declare(strict_types=1);

namespace App\MerchantReview\Domain;

final class MerchantReviewCreated implements Event
{
    private function __construct(
        private MerchantReview $merchantReview
    ) {}

    public static function create(
        MerchantReview $merchantReview
    ): self {
        return new self($merchantReview);
    }

    public function getReview(): MerchantReview
    {
        return $this->merchantReview;
    }
}
